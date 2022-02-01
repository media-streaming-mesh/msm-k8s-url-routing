package url_handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	ClusterIp net.IP
	Endpoints []net.IP
}

type UrlHandler struct {
	clientset *kubernetes.Clientset
}

func (uh *UrlHandler) NewUrlHandler() {
	u := &UrlHandler{}

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("error creating out of cluster config: %v", err)
	}

	u.clientset = clientset
}

func (uh *UrlHandler) log(format string, args ...interface{}) {
	// keep remote address outside format, since it can contain %
	log.Println("[K8S API Client] " + fmt.Sprintf(format, args...))
}

func (uh *UrlHandler) GetUrls(Url string) []string {
	uh.log("url provided: %s", Url)
	u, err := url.Parse(Url)

	if err != nil {
		empty := make([]string, 1)
		uh.log("could not parse url: %v", err)
		return empty
	}

	// lookup IP and port
	endpoints := uh.resolveEndpoints(u.Hostname(), u.Port(), u.Path)
	urls := make([]string, len(endpoints))

	for i, endpoint := range endpoints {
		cpy := u
		cpy.Host = endpoint
		urls[i] = cpy.String()
	}
	return urls
}

func (uh *UrlHandler) resolveEndpoints(hostname string, port string, path string) []string {
	var (
		serviceName string
		addresses   []string
	)

	uh.log("get resources for %s, %s, %s", hostname, port, path)

	if hostname == "localhost" {
		if len(path) > 0 {
			uh.log("IP is localhost - path is %s", path)
			serviceName = path[1:]
		}
	} else {
		// look up hostname in DNS if needed
		addresses = uh.resolveHost(hostname)

		// if DNS resolved to set of addrs, return addrs
		if len(addresses) > 1 {
			uh.log("DNS returned %d results", len(addresses))
			endpoints := make([]string, len(addresses))
			for i := range addresses {
				endpoints[i] = addresses[i] + ":" + port + path
			}
			return endpoints
		}

		// see if IP is a node IP
		if uh.isNodeIP(addresses[0]) {
			uh.log("IP matches node IP")
			if len(path) > 0 {
				uh.log("path is %s", path)
				serviceName = path[1:]
			}
		}
	}

	// else if single ip is a clusterIP, return endpoints
	if serviceName == "" {
		var err error
		serviceName, err = uh.getServiceName(addresses[0]) // addresses is a slice of length 1

		// if ip is not clusterIP, return ip
		if err != nil {
			uh.log("unable to look up cluster IP, returning IP given (%s)", addresses[0])
			return []string{addresses[0] + ":" + port + "/" + path}
		}

		uh.log("ClusterIP's service name is %s", serviceName)
	}

	// search for endpoints that belong to this service
	uh.log("Path or clusterIP was resolved to the service: %s", serviceName)
	return uh.getEndpoints(serviceName)
}

// get all endpoints for a named service
func (uh *UrlHandler) getEndpoints(serviceName string) []string {
	ends, err := uh.clientset.CoreV1().Endpoints("default").Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		uh.log("failed to get service endpoints")
		return []string{}
	}

	endpoints := make([]string, 0)
	for _, end := range ends.Subsets {
		addrs := end.Addresses
		ports := end.Ports

		for _, addr := range addrs {
			for _, port := range ports {
				uh.log("service endpoint: %s:%d", addr.IP, port.Port)
				endpoint := addr.IP + ":" + strconv.FormatUint(uint64(port.Port), 10)
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	return endpoints
}

// check if IP is a node IP
func (uh *UrlHandler) isNodeIP(hostname string) bool {
	// get node IPs
	nodes, err := uh.clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		uh.log("failed to get node list")
		return false
	}

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Address == hostname {
				uh.log("matched node address")
				return true
			}
		}
	}

	return false
}

// look for a service with clusterIP matching a given IP
func (uh *UrlHandler) getServiceName(clusterIP string) (string, error) {
	services, err := uh.clientset.CoreV1().Services("default").List(context.TODO(),
		v1.ListOptions{})

	if err != nil {
		uh.log("unable to get list of services")
		return "", nil
	}

	for _, service := range services.Items {
		if service.Spec.ClusterIP == clusterIP {
			uh.log("found service name %s for clusterIP %s", service.Name, clusterIP)
			return service.Name, nil
		}
	}

	return "", fmt.Errorf("service name could not be resolved")
}

// if given hostname is an IP return it, else perform DNS lookup
// DNS may return multiple IPs (uh.g. headless services)
func (uh *UrlHandler) resolveHost(host string) []string {
	var ip net.IP

	uh.log("host = %s", host)

	ip = net.ParseIP(host) // returns nil if not ip

	// if host is not an ip, perform dns lookup
	if ip == nil {
		uh.log("looking up in DNS")
		addrs, err := net.LookupHost(host)
		if err != nil {
			uh.log("could not resolve resource to ips: %s, returning given resource: %s", err.Error(), host)
			return []string{host}
		}
		uh.log("resolution: ", addrs)
		return addrs
	}
	// dns not needed, already an ip
	uh.log("resolution: ", ip)
	return []string{ip.String()}
}
