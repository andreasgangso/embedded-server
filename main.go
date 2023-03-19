package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/ncruces/zenity"
)

func main() {

	// Serve static files using go.rice
	httpBox := rice.MustFindBox("static")
	httpBoxHandler := http.FileServer(httpBox.HTTPBox())

	// Prompt the user for the port number
	port := promptForPort("Enter a http port (1-65535).", "2020")
	tlsPort := promptForPort("Enter a https port (1-65535).", "2021")
	addr := fmt.Sprintf(":%d", *port)
	tlsAddr := fmt.Sprintf(":%d", *tlsPort)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: httpBoxHandler,
	}
	tlsServer := &http.Server{
		Addr:    tlsAddr,
		Handler: httpBoxHandler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{generateSelfSignedCert()},
		},
	}

	errs := make(chan error)

	// Starting HTTP server
	go func() {
		fmt.Printf("Starting server on http://localhost%s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errs <- err
		}
	}()

	// Starting HTTPS server
	go func() {
		fmt.Printf("Starting server on https://localhost%s\n", tlsAddr)
		if err := tlsServer.ListenAndServeTLS("", ""); err != nil {
			errs <- err
		}
	}()

	go func() {
		openInBrowser(fmt.Sprintf("http://localhost%s/", addr))
		zenity.Info(fmt.Sprintf("The app is ready at http://localhost%s and https://localhost%s.", addr, tlsAddr), zenity.OKLabel("Stop the app"))
		os.Exit(0)
	}()

	select {
	case err := <-errs:
		zenity.Error(fmt.Sprintf("Could not start serving service due to (error: %s)", err), zenity.OKLabel("OK"))
		log.Printf("Could not start serving service due to (error: %s)", err)
	}

}

func promptForPort(title string, defaultValue string) *int {
	port, err := zenity.Entry(title, zenity.Title("Enter port"), zenity.EntryText(defaultValue))
	if err != nil {
		log.Fatal("Error getting port number:", err)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		zenity.Error("Invalid port number. Please enter a number between 1 and 65535.")
		return nil
	}
	return &portNumber
}

func openInBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unknown platform %s", runtime.GOOS)
	}
	if err != nil {
		log.Printf("Error opening browser: %s", err)
	}
}

func generateSelfSignedCert() tls.Certificate {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Example Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, _ := tls.X509KeyPair(certPem, keyPem)
	return cert
}
