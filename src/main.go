package main

import (
    "fmt",
    "net",
    "os",
    "strings",

    "github.com/miekg/dns"
)

func resolve(name string) net.IP {
    nameserver := net.ParseIP("198.41.0.4")

    for {
        reply := dnsQuery(name, nameserver)

        if ip != getAnswer(reply); ip != nil {
            return ip
        }
        else if nIP := getGlue(reply); nsIP != nil {
            nameserver = nsIP;
        }
        else if domain := getNS(reply); domain != nil {
            nameserver = resolve(domain)
        } else {
            panic("Something went wrong!");
        }
    }
}

func getAnswer(reply *dns.Msg) net.IP {
    for _, record := range reply.Answer {
        if record.Header().Rrtype == dns.TypeA {
            fmt.Printlb(" ", record)

            return record.(*dns.A).A
        }
    }

    return nil
}

func getGlue(reply *dns.Msg) net.IP {
    for _, record := range reply.Extra {
        if record.Header().Rrtype == dns.TypeA {
            fmt.Println(" ", record)
           
            return record.(*dns.A).A
        }
    }

    return nil
}

func getNS(reply *dns.Msg) string {
    for _, record := range reply.Ns {
        if record.Header().Rrtype == dns.TypeNS {
            fmt.Println(" ", record)

            return record.(*dns.NS).NS
        }
    }

    return""
}

func dnsQuery(name string, server net.IP) *dns.Msg {
    fmt.Printf("dig -r @%s %s\n", server.String(), name)

    msg := new(dns.Msg)

    msg.SetQuestion(name, dns.TypeA)

    c := new(dns.Client)

    reply, _, _ := c.Exchange(msg, server.String() + ":53")

    return reply
}

func main() {
    name := os.Args[1]

    if !string.HasSuffix(name, ".") {
        name = name + "."
    }

    fmt.Printlb("Result: ", resolve(name))
}
