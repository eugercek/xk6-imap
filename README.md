# Examples

## Gmail

Google removed plain username+password authentication for 3rd party. You need to create a "App password" to be able to use username+password for authentication.

```js
import Imap from "k6/x/imap";

export default function () {
  const [message, error] = Imap.read(
    "my_email@gmail.com",
    "password123",
    "imap.gmail.com",
    993,
    {
      Subject: ["Verify your email"],
    }
  );

  if (error != "") {
    console.error(error);
  } else {
    console.log(message);
  }
}
```

## `emailClient`

Use email client if you need to read multiple messages and don't want to login everytime.

```js
import imap from "k6/x/imap";

export default function () {
  const client = imap.emailClient(
    "my_email@gmail.com",
    "password123",
    "imap.gmail.com",
    993
  );

  const loginError = client.login();

  if (loginError != "") {
    console.error(loginError);
    return;
  }

  let [message, err] = client.read({
    Subject: ["Verify your email"],
  });

  if (err != "") {
    console.error(err);
    return;
  }

  console.log(message);
  client.logOut();
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
