package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dasio/base45"
	"github.com/fxamacker/cbor/v2"
)

func main() {
	var encodedGarbage string = ""

	// file, err := os.Open(os.Stdin)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		encodedGarbage += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Scanning stdin failed:", err)
	}
	// fmt.Println("Garbage out:", encodedGarbage)

	// According to EU specs, the first three characters identify the version of the certificate.
	// See @https://ec.europa.eu/health/sites/default/files/ehealth/docs/digital-green-certificates_v3_en.pdf
	// Relevant for us are versions HC1, HC2, etc. (HC = Health Certificate)

	decodedGarbage, err := base45.DecodeString(encodedGarbage[4:])
	if (err != nil) {
		fmt.Println("Could not decode QR garbage; error was:", err)
	} else {
		// fmt.Printf("QR data decoded with Base45: %#v\n", decodedGarbage)
	}
	// fmt.Printf("Forcing it to be a string: %s\n", string(decodedGarbage[4:]))

	// The same specs assert that the data is now zlib'ed, so we must de-zlib it:
	dezlibedGarbageBuffer := bytes.NewReader(decodedGarbage)

	dezlibedGarbage, err := zlib.NewReader(dezlibedGarbageBuffer)
	if err != nil {
		log.Fatal("Could not create a handler for the de-zlib garbage:", err)
	}

	ucviPayload, err := io.ReadAll(dezlibedGarbage)
	if err != nil {
		log.Fatal("De-zlib'ing failed utterly:", err)
	}
	fmt.Printf("UCVI Payload is now: %q\n", ucviPayload);

	// now unmarshal the de-zlibed payload as CBOR.
	// Using https://github.com/fxamacker/cbor#usage

	// Part of COSE header definition
	type coseHeader struct {
		Alg int    `cbor:"1,keyasint,omitempty"`	 // Cryptographic algorithm used, ES256 or PS256
		Kid []byte `cbor:"4,keyasint,omitempty"`	 // Key Identifier (first 8 bytes of the hash value)
		//		IV  []byte `cbor:"5,keyasint,omitempty"`	// Initialisation Vector
	}

	type commonPayloadValues struct {
		Iss	   []byte	 `cbor:"2,omitempty"`	 // Issuer of the Digital COVID-19 Certificate (DGC)
		Iat	   []byte	 `cbor:"2,omitempty"`	 // Issuing Date of the DGC
		Exp	   []byte	 `cbor:"2,omitempty"`	 // Expiring Date of the DGC
		Hcert []byte `cbor:"5,omitempty"`	 // Payload of the DGC (Vac[cined], Tst[result], Rec[overed])
	}

	// Signed CWT is defined in RFC 8392
	type signedCWT struct {
		_           struct{} `cbor:",toarray"`
		Protected   []byte
		Unprotected coseHeader	 `cbor:",toarray"`
		// Protected   coseHeader	`cbor:",toarray"`
		// Payload     commonPayloadValues	`cbor:",toarray"`
		Payload	   []byte
		Signature  []byte
		// Unprotected coseHeader	`cbor:",omitempty"`	// Digital COVID-19 Certificates do not have unprotected data as per specs
	}

	var cwt signedCWT

	// Unmarshal CBOR (binary json) data structure
	// data is []byte containing signed CWT
	if err := cbor.Unmarshal(ucviPayload, &cwt); err != nil {
		log.Panic("CBOR unmarshal failed:", err)
	}

	fmt.Printf("Unmarshalled CBOR CWT %+v\n", cwt)
	// if this worked correctly, we should now be able to get JSON out of this?
	//var e interface{}
	//err = json.Unmarshal(emptyInterface, &e)

	fmt.Println("\n---")
	//fmt.Printf("Protected: %s\n", cwt.Protected)
	fmt.Printf("Algorithm field (unprotected): %d\n", cwt.Unprotected.Alg)
	fmt.Printf("Key Identifier  (unprotected): %d\n", cwt.Unprotected.Kid)
	// fmt.Printf("Initialisation Vector (unprotected): %d\n", cwt.Unprotected.IV)
	fmt.Printf("Protected: %v\n", cwt.Protected)
	fmt.Printf("Last 8 bytes: %s\n", cwt.Protected[5:])
	fmt.Printf("Payload: %s\n", cwt.Payload)
	// fmt.Printf("Payload - Issuer: %q\n", cwt.Payload.Iss)
	// fmt.Printf("Payload - Issuing date: %q\n", cwt.Payload.Iat)
	// fmt.Printf("Payload - Expiration date: %q\n", cwt.Payload.Exp)
	// fmt.Printf("Payload - Health certificate: %q\n", cwt.Payload.Hcert)
	fmt.Printf("Signature: %v\n", cwt.Signature)

	// Attempt to Unmarshall the Payload as JSON...
	fmt.Println("\n---\nHealth Certificate as follows:\n")

	var e interface{}
	// if err := json.Unmarshal(cwt.Payload, &e); err != nil {
	// 	log.Panic("JSON unmarshalling of CWT health certificate payload failed:", err)
	// }
	var payload commonPayloadValues
	if err := cbor.Unmarshal(cwt.Payload, &payload); err != nil {
		log.Panic("CBOR unmarshalling of CWT health certificate payload failed:", err)
	} else {
		fmt.Printf("Parsed payload: %v\n", payload)
	}

	// Use type assertion to access underlying map[string]interface{}
	m := e.(map[string]interface{})

	// Use a type switch to access values as their concrete types
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is", vv)
		case float64:
			fmt.Println(k, "is", vv)
		case []interface{}:
			fmt.Println(k, ":")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is an unknown type")
		}
	}
}
