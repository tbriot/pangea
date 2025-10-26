wrk.method = "POST"

wrk.headers["Accept"] = "*/*"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "GitHub-Hookshot/e228741"
wrk.headers["X-GitHub-Delivery"] = "8c026ab6-b283-11f0-8604-75e0a62e4373"
wrk.headers["X-GitHub-Event"] = "push"
wrk.headers["X-GitHub-Hook-ID"] = "504572572"
wrk.headers["X-GitHub-Hook-Installation-Target-ID"] = "1011602"
wrk.headers["X-GitHub-Hook-Installation-Target-Type"] = "integration"
wrk.headers["X-Hub-Signature"] = "sha1=560dc6f612e008bd6b470142d0ec442d6b4450be"
wrk.headers["X-Hub-Signature-256"] = "sha256=cdb4c4f3e2afd6599a4591037eb2e7ad59583c1f1fbb08ad03b73d5cbd9d954e"

file = io.open("./test/01_push_event_payload.json", "rb")
wrk.body = file:read("*a")