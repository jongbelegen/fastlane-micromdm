package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	bolt "go.etcd.io/bbolt"
)

func main() {
	dbPath := "micromdm.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	type certInfo struct {
		Key       string
		NotBefore time.Time
		NotAfter  time.Time
	}

	var certs []certInfo

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scep_certificates"))
		if b == nil {
			log.Fatal("scep_certificates bucket not found")
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if key == "ca_certificate" {
				cert, err := x509.ParseCertificate(v)
				if err == nil {
					expired := ""
					if time.Now().After(cert.NotAfter) {
						expired = " ** EXPIRED **"
					}
					fmt.Printf("=== CA Certificate ===\n")
					fmt.Printf("  Subject:    %s\n", cert.Subject)
					fmt.Printf("  NotBefore:  %s\n", cert.NotBefore.Format("2006-01-02"))
					fmt.Printf("  NotAfter:   %s\n", cert.NotAfter.Format("2006-01-02"))
					fmt.Printf("  Status:     %.1f years remaining%s\n\n", time.Until(cert.NotAfter).Hours()/24/365.25, expired)
				}
				return nil
			}
			if key == "ca_key" || key == "serial" {
				return nil
			}
			cert, err := x509.ParseCertificate(v)
			if err != nil {
				return nil
			}
			certs = append(certs, certInfo{
				Key:       key,
				NotBefore: cert.NotBefore,
				NotAfter:  cert.NotAfter,
			})
			return nil
		})
		return nil
	})

	sort.Slice(certs, func(i, j int) bool {
		return certs[i].NotBefore.After(certs[j].NotBefore)
	})

	fmt.Printf("=== Device Certificates (newest first) ===\n\n")
	limit := 30
	if len(os.Args) > 2 && os.Args[2] == "--all" {
		limit = len(certs)
	}
	if len(certs) < limit {
		limit = len(certs)
	}

	for _, c := range certs[:limit] {
		years := c.NotAfter.Sub(c.NotBefore).Hours() / 24 / 365.25
		status := "OK"
		if time.Now().After(c.NotAfter) {
			status = "EXPIRED"
		}
		fmt.Printf("  %-50s  %s → %s  (%5.1f yrs)  [%s]\n",
			c.Key, c.NotBefore.Format("2006-01-02"), c.NotAfter.Format("2006-01-02"), years, status)
	}

	fmt.Printf("\n  Total: %d certificates\n", len(certs))
}
