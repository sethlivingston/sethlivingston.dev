---
layout: "../../layouts/BlogPage.astro"
title: "A complete guide to Fetch API error handling"
description: "Sending HTTP requests with the Fetch API is easy—if everything goes well. But if something goes wrong, the Fetch API turns into a mess. This guide shows you how to untangle the mess and add robust error handling."
pubDate: "29 Oct 2022"
draft: true
---

There are four types of errors that can happen when you send HTTP requests with the [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API).

<dl>
<dt>Pre-processing error</dt>
<dd>Happens when the request is malformed</dd>
<dt>Network error</dt>
<dd>Happens when the request can't reach the remote server</dd>
<dt>HTTP error</dt>
<dd>Happens when the remote server responds with an error condition</dd>
<dt>Post-processing error</dt>
<dd>Happens when the response is malformed</dd>
</dl>

For each type of error, you need to know:

1. What causes it to happen
1. How often should you expect to see it in production
1. How to detect it
1. How to handle it
1. How to test it

Let's start with the error you should see the least often.

## Pre-processing errors

**Cause**: Pre-processing errors happen when your request violates 1) HTTP protocol rules or 2) Fetch API parameter rules. Example violations include:

* Invalid URL protocol, like `foo://example.com`
* Including a body on a `GET` request
* Invalid characters in header names, like `"no spaces please": "value"`

**How often**: You should detect and fix pre-processing errors in development. They are preventable, so they should be rare in production. 

**How to detect**: The `fetch()` function throws a `TypeError` exception when it detects a pre-processing error. The [MDN docs](https://developer.mozilla.org/en-US/docs/Web/API/fetch#exceptions) include a complete list of possible pre-processing errors. 

**How to handle**: Catch `TypeError` exceptions when you call the fetch function. Log the message so you know the details of the violation. Provide users with guidance if you can, but this is probably an error only a developer can fix.

```javascript
let response, error;

try {
  response = await fetch(resource, options);
  // No pre-processing errors

} catch (err) {
  if (err.constructor === TypeError) {
    // Pre-processing error, log so developers can troubleshoot and fix.
    // The user is likely stuck.
    console.error(err.message);
  }
}
```

**How to test**: Use the facilities in your testing tools to mock or intercept the `fetch()` function and throw an exception. 

There's more to say about pre-processing errors in the next section.

## Network errors

**Cause**: Network errors happen when anything prevents `fetch()` from making a connection to the remote server. Here are just a few ways that can happen:

* Wifi disconnects
* Router or local network goes down
* DNS can't resolve the host name

I also group abort errors into the category of network errors. Abort errors happen when the request is aborted due to a call to `AbortController.abort()` ([see docs](https://developer.mozilla.org/en-US/docs/Web/API/AbortController/abort)).

**How often**: Network errors happen frequently, and there's nothing you can do to prevent them.

**How to detect**: Here's where things get messy. Unfortunately `fetch()` signals network errors the same way it signals pre-processing errors: it throws a `TypeError` exception.

```javascript
let response, error;

try {
  response = await fetch(resource, options);
  // No network errors OR pre-processing errors

} catch (err) {
  if (err.constructor === TypeError) {
    // Network error OR pre-processing error, log so developers
    // can troubleshoot and fix.
    console.error(err.message);
  }
}
```

Could we use the exception's `.message` property to differentiate between the two types of exceptions? I recommend against this for two reasons. First, using human-readable strings to make processing decisions leads to fragile, unreliable code. Second, the `.message` property value isn't standardized.

The table below shows the `TypeError.message` value for various network and pre-processing errors. These results are from an ad-hoc test I ran on MacOS[^1].

| Error condition               | Chrome's message                                                                      | Safari's message                                | Node.js's message                                        |
| ----------------------------- | ------------------------------------------------------------------------------------- | ----------------------------------------------- | -------------------------------------------------------- |
| No network                    | Failed to fetch                                                                       | Load failed                                     | fetch failed                                             |
| Invalid URL                   | Failed to fetch                                                                       | Load failed                                     | fetch failed                                             |
| Sending body with GET request | Failed to execute 'fetch' on 'Window': Request with GET/HEAD method cannot have body. | Request has method 'GET' and cannot have a body | Request with GET/HEAD method cannot have body.           |
| Invalid header name           | Failed to execute 'fetch' on 'Window': Invalid name                                   | Invalid header name: 'sp ace'                   | Headers.append: "sp&nbsp;ace" is an invalid header name. |

There is no reliable way to differentiate between pre-processing errors and network errors.

Because pre-processing errors are rare and largely preventable, it's acceptable to treat them as network errors. In this way network errors become "any error that prevents fetch from making a connection." What's important is that you log the exception message so developers can spot pre-processing errors sooner than later. If you run thorough tests on your application or library, then you can assume exceptions in production are true network errors.

**How to test**: Use the facilities in your testing tools to mock or intercept the `fetch()` function and throw an exception.

For manual tests, you can turn off or throttle networking three ways: 

1. Operation system settings
2. Web browser's dev console
3. Using a network analyzer like [Wireshark](https://www.wireshark.org)

Handling network errors effectively is the sign of a robust application or library. Your users will thank you. Okay, they probably won't thank you. But at least they'll complain less.

### HTTP errors

### Post-processing errors

## Should you wrap fetch yourself or use an existing Fetch API wrapper?



[^1]: You can try this yourself. Start with this function. Make modifications as needed to simulate different error conditions.
      ```javascript
      (async function() {
        try { 
          await fetch("https://api.chucknorrisnopenopenope.io/jokes/random"); 
        } catch (err) { 
          console.log("message:", err.message); 
        } 
      })()
      ```