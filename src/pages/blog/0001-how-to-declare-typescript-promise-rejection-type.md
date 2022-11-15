---
layout: "../../layouts/BlogPage.astro"
title: "The rejection type for a Promise is always `any`"
description: "The template parameter for Promise<T> is the fulfillment type — not the rejection type. It isn't possible to declare the rejection type."
pubDate: "8 Nov 2022"
---

Consider a function that returns a rejected Promise. It calls `Promise.reject` with a custom object of type `Details`.

```typescript
export interface Details {
  message: string;
}

export function rejectWithDetails(msg: string): Promise<Details> {
                                                ^^^^^^^^^^^^^^^^
  return Promise.reject({ message: msg });
}
```

This compiles, but the return type for `rejectWithDetails` is incorrect. The return type declaration says, "This function returns a promise that resolves to a `Details` object when fulfilled." That's not what we want. We're trying to convey the rejection type.

We can't declare the rejection type because the `Promise` type doesn't have a template parameter for it. In `Promise<T>`, `T` is the **fulfillment** type. The **rejection** type is always `any`. 

The correct declaration for the function is:

```typescript
export function rejectWithDetails(msg: string): Promise<any> { ...
                                                ^^^^^^^^^^^^
```

This is still a little confusing, though. This function does not return a `Promise<any>` because the rejection type is always `any`. It returns a `Promise<any>` because it doesn't know what the Promise's fulfilled type is. Only the calling function knows.

It is not possible to declare a Promise's rejection type. It is always `any`.

### Sources

* [Different types for rejected/fulfilled Promise · Issue #7588 · microsoft/TypeScript](https://github.com/Microsoft/TypeScript/issues/7588)
* [javascript - Typescript Promise rejection type - Stack Overflow](https://stackoverflow.com/questions/50071115/typescript-promise-rejection-type/50071254#50071254)



