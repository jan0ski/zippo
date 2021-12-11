package main

import (
	zippo "github.com/jan0ski/zippo/pkg"
)

func main() {
	s := &zippo.Server{
		Config: &zippo.ServerConfig{
			ButaneTemplate: "./config.yaml",
			CNIVersion: "v0.8.2",
			CRICtlVersion: "v1.17.0",
			K8sVersion: "v1.20.0",
			SSHUser: "core",
			SSHPubkey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCtX+XaqAX5Xa58ZLteuvKFA5NhylWlLDngRYcPvVWVyxDUrA2w7F+ovH1tdT8TNIZYBQPWN7ioOFGKPwOmI0wCxbqA/X2ZAXzt6Icq/ELqD+w7oaMdZQIHPWcOz2kxA1G92Q/MvUOMmIz8GgeCAWugB5lwTd3iWBS9FxWIk4Z3hVCSKtTDpqn7VJpF3wGRN3ZoYRdvi5U+q36kjz7k17yyYQcLoHm5HgH1Wj5MFi3JidS8bF7t2lHoyi2t/R6gBXdASUr6vCxbBbokFHEhq7VlmmwInuva77Qr/StGNfjEwpva7VXpvRnRdHk2cpioBAXbY+XJH2dMBWhBmCqktxTzSgnqhwt0FO6ph3ukYeOdxr3x9CVkWm9UeqpgPchCtNPYlQjxFRCEM18JmuH77jXqwYC7n2tshQJ1xr9PJRSSaO4EAbdxJoloLX/tIl1dkGLMlxGzrfedwqru82L3qbqxOwp90Mdf8r4OvhUhfD6/iAKHiAXiaameqz1N+BDj2RIbjxMTWnmTk2CwuG7A0FHW6TkSbgVvJUfq+eTMzIbMVOIjPNF8tjXBksJeoA4mi8bhLpC814pvYLqYRDImrRRv8WDMID/djEC8ajqLnCA3CHxZNXenm0A9gJsQrs2kF6xo/p62GAKnNxmE7sMdxV/mwHUVDSkPdREsPC6h/ID0uQ== jan0ski@Zeus",
		},
	}
	s.Run()
}
