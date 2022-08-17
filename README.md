# Examples

## Read Gmail mail

Google removed plain username+password for 3rd party. You need to create a "App password" to be able to use username+password for authentication.

```js
import Imap from "k6/x/imap";

export default function () {
  const message = Imap.read(
    "my_email@gmail.com",
    "password123",
    "imap.gmail.com",
    993,
    {
      Subject: ["Verify your email"],
    }
  );
}
```

# Build

Don't forget to use this binary instead of the `k6` binary in your path.

```bash
xk6 build --with github.com/eugercek/xk6-imap

# ./k6 run script.js
```

# TODO List

- Give examples for major email providers
- Give examples for how to measure elapsed time
- Create unit tests for the Go code
- Expose more query options
- Investigate error handling on xk6 extensions
