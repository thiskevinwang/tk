package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Session struct {
	Id    string `json:"sessionId"`
	Key   string `json:"sessionKey"`
	Token string `json:"sessionToken"`
}

type TokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

var awsCommand = &cobra.Command{
	Use:   "aws",
	Short: "launch the aws console in your browser",
	Long:  "launch the aws console in your browser",
	Run: func(cmd *cobra.Command, args []string) {
		region := "us-east-1"
		console_url := "https://console.aws.amazon.com/?region=" + region
		sign_in_url := "https://signin.aws.amazon.com/federation"

		name := "stranger_danger"
		if len(args) >= 1 {
			name = args[0]
		}

		ctx := context.TODO()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatal(err)
		}

		client := sts.NewFromConfig(cfg)

		creds, err := client.GetFederationToken(ctx, &sts.GetFederationTokenInput{
			Name: aws.String(name),

			PolicyArns: []types.PolicyDescriptorType{
				types.PolicyDescriptorType{
					Arn: aws.String("arn:aws:iam::aws:policy/ReadOnlyAccess"),
				},
				types.PolicyDescriptorType{
					Arn: aws.String("arn:aws:iam::aws:policy/AWSBillingReadOnlyAccess"),
				},
			},
		})

		if err != nil {
			log.Fatal(aurora.Red(err))
		}

		baseUrl, err := url.Parse(sign_in_url)
		if err != nil {
			log.Fatal(aurora.Red("Malformed URL: " + err.Error()))
			return
		}

		session := Session{
			Id:    aws.ToString(creds.Credentials.AccessKeyId),
			Key:   aws.ToString(creds.Credentials.SecretAccessKey),
			Token: aws.ToString(creds.Credentials.SessionToken),
		}
		b, err := json.Marshal(session)
		if err != nil {
			log.Fatal(aurora.Red(err))
			return
		}

		params := url.Values{}
		params.Add("Action", "getSigninToken")
		params.Add("SessionDuration", "43200")
		params.Add("Session", string(b))

		baseUrl.RawQuery = params.Encode()

		uri := baseUrl.String()

		resp, err := http.Get(uri)
		if err != nil {
			log.Fatal(aurora.Red(err))
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(aurora.Red(err))
			}

			data := TokenResponse{}
			json.Unmarshal(bodyBytes, &data)

			params := url.Values{}
			params.Add("Action", "login")
			params.Add("Issuer", "arbitrary_issuer")
			params.Add("Destination", console_url)
			params.Add("SigninToken", data.SigninToken)

			baseUri, err := url.Parse(sign_in_url)
			if err != nil {
				log.Fatal(aurora.Red("Malformed URL: " + err.Error()))
				return
			}

			baseUri.RawQuery = params.Encode()
			console_uri := baseUri.String()

			log.Println("Attempting to open AWS console for:", aurora.Blue(region))

			openbrowser(console_uri)
		} else {
			log.Warn(aurora.Yellow("Non 200 status code"))
		}
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(awsCommand)
}
