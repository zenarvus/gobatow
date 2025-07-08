# gobatow
Go back to work!

Block websites until your required md-agenda tasks are completed.

## Quick Start
1. Add the `#blck` tag to tasks to block access to websites until they are completed.
    - This app uses the scheduled time property. If it does not exist, it blocks access until the task is no longer a TODO or HABIT. If it exists, access is blocked if the current time exceeds the scheduled time.

2. Clone the repository and navigate into it:
   ```bash
   git clone https://github.com/zenarvus/gobatow && cd gobatow
   ```

3. Create a `config.go` file with the following content and customize it according to your needs:
   ```go
   package main

   var proxyPort = "8383"
   var queryPort = "8081"

   // For the query server; ignored if empty.
   const certPath = "cert.pem"
   const keyPath = "key.pem"

   // 24-hour time format HH:MM
   // Ignored if an empty string
   var allowBefore = "06:30" // Allow access to blocked sites before this time. Useful for your tasks that requires access to the blocked websites before this time.
   var blockAfter = "23:00" // Block access to sites after this time to prevent disrupting your sleep pattern.

   // Local agenda file paths. Consider using Syncthing or a similar tool if you plan to run this on an external server.
   var agendaFiles = []string{
       "/home/user/agenda/habits.md",
   }

   const blockType = "blacklist" // Options: "blacklist" or "whitelist"

   // Domains that contain these strings are blocked.
   var blockedDomains = []string{
       "reddit",
       "youtube",
       "mastodon",
   }
   ```

4. Configure your device to connect to the proxy server.

5. Start the server:
   ```bash
   go run *.go
   ``` 
> You can also make GET requests to the root path of the query port, which returns "tasks-completed" or "tasks-uncompleted" in plain text.
> - The returned value can be processed by a client-side browser script or a similar method to block access to certain sites.
