package entity

var ImportantWindowsPorts = []struct {
	port        int
	service     string
	protocol    string
	description string
}{
	{
		port:        1337,
		service:     "砳 Edge Dev Tools",
		protocol:    "TCP",
		description: "Used by Microsoft Edge Developer Tools for remote debugging.",
	},
	{
		port:        1433,
		service:     "Microsoft SQL Server",
		protocol:    "TCP",
		description: "Default port for Microsoft SQL Server database engine.",
	},
	{
		port:        1434,
		service:     "Microsoft SQL Monitor",
		protocol:    "TCP/UDP",
		description: "Listens for incoming connections to SQL Server instances that are not listening on the default port.",
	},
	{
		port:        1723,
		service:     "Point-to-Point Tunneling Protocol (PPTP)",
		protocol:    "TCP",
		description: "Used for implementing virtual private networks (VPNs).",
	},
	{
		port:        1801,
		service:     "Microsoft Message Queuing (MSMQ)",
		protocol:    "TCP/UDP",
		description: "Enables applications to communicate with each other across a network.",
	},
	{
		port:        2382,
		service:     "Microsoft SQL Server Analysis Services (SSAS)",
		protocol:    "TCP",
		description: "Default port for SQL Server Analysis Services.",
	},
	{
		port:        2383,
		service:     "Microsoft SQL Server Analysis Services (SSAS)",
		protocol:    "TCP",
		description: "Used for browser requests and other instances of Analysis Services.",
	},
	{
		port:        3268,
		service:     "Microsoft Global Catalog",
		protocol:    "TCP",
		description: "Used to search for objects in an Active Directory forest.",
	},
	{
		port:        3269,
		service:     "Microsoft Global Catalog (SSL)",
		protocol:    "TCP",
		description: "Secure version of the Global Catalog service using SSL.",
	},
	{
		port:        3389,
		service:     "Remote Desktop Protocol (RDP)",
		protocol:    "TCP",
		description: "Default port for Remote Desktop services, allowing remote access to a Windows machine.",
	},
	{
		port:        4600,
		service:     "Microsoft",
		protocol:    "TCP",
		description: "Part of the Azure Communication Services for real-time communication.",
	},
	{
		port:        5009,
		service:     "WinFS",
		protocol:    "TCP/UDP",
		description: "Associated with the cancelled Windows Future Storage project, but may still be reserved.",
	},
	{
		port:        5357,
		service:     "Web Services for Devices (WSDAPI)",
		protocol:    "TCP/UDP",
		description: "Used for device discovery and services on a local network.",
	},
	{
		port:        5722,
		service:     "Microsoft DFS Replication Service",
		protocol:    "TCP",
		description: "Used by the Distributed File System Replication service.",
	},
	{
		port:        5985,
		service:     "Windows Remote Management (WinRM)",
		protocol:    "TCP",
		description: "Default port for WinRM over HTTP, used for remote management of Windows servers.",
	},
	{
		port:        5986,
		service:     "Windows Remote Management (WinRM) over HTTPS",
		protocol:    "TCP",
		description: "Secure version of WinRM using HTTPS.",
	},
	{
		port:        8080,
		service:     "HTTP Alternate",
		protocol:    "TCP",
		description: "Commonly used as an alternative port for web servers, and may be used by various applications.",
	},
	{
		port:        9001,
		service:     "Tor (unofficial)",
		protocol:    "TCP",
		description: "Often used by the Tor service for its control port.",
	},
	{
		port:        9389,
		service:     "Active Directory Web Services",
		protocol:    "TCP",
		description: "Provides a web service interface to Active Directory.",
	},
}

func IsImportantWindowsPort(port int) bool {
	for _, p := range ImportantWindowsPorts {
		if p.port == port {
			return true
		}
	}
	return false
}
