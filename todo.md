# Next Steps:

## v-for:
* ~~for is broken have to check the index var too and yell if its empty~~
* ~~and why didn't i save the iterable in loopinfo ????~~
## more specififity in error messages
### handle the parts of for directive separately

**missing expression after 'in'**

**no "in"**

**"," without preceeding "key" or successive "var"**
## have to disallow any if and for on same element
so v-for and v-if might create racing condition because of evaluation order (when is the loop variable ready?)
encourage use of virtual tag (template) to specify order of precedence (enclosing v-for or v-if)
## must still handle components  - can change syntax as long as vue/vscode tolerates syntax
## virtual tags as template definitions  - can change syntax as long as vue/vscode tolerates syntax
## slots  - can change syntax as long as vue/vscode tolerates syntax
## props  - can change syntax as long as vue/vscode tolerates syntax
## attribute bindings
## aliases v-bind, v-on
## specialty class bindings - can change syntax as long as vue/vscode tolerates syntax
## very special style bindings - can change syntax as long as vue/vscode tolerates syntax
## 

# `<template>` tag directives for virtual element in Vue:

**OK, here's the list of directives that a `<template>` tag can accept when used as a virtual element in Vue.js:**
* v-if: Conditionally renders the template's content based on a boolean expression.
* v-else: Used along with v-if to render alternative content if the v-if condition is false.
* v-else-if: Chained with v-if to provide additional conditional branches.
* v-for: Repeats the template's content for each item in an array.
* v-show: Conditionally shows or hides the template's content (using CSS display property).
* v-once: Renders the template's content only once and then caches it, preventing further updates.
* v-pre: Skips compilation for the template's content, rendering it as raw HTML.
* v-cloak: Hides the template's content until the Vue instance is fully compiled (useful for preventing flickering).

### v-cloak
maybe use SSR to adress this

**Important notes**

The v-for directive is particularly useful with virtual `<template>` tags to create loops without introducing extra DOM elements.
The :key attribute is important when using v-for to help Vue identify and track elements in the loop efficiently.
Virtual `<template>` tags do not render any output themselves; they act as containers for the directives and their associated content.

## v-show
still has to be implemented

## v-once
we need to look at SSR usage this is a prime candidate being useful for SSR

## v-pre
must absolutely implement this it is basically rendering the template as is and treats all directives and
interpolation as text

e.g.

```vue
<template>
  <div v-pre>{{ message }}</div>
</template>
...
```

```
Hello, {{ name }}!
```

Without `v-pre`, Vue would try to interpret {{ name }} as a data binding expression, which might result in an error or unexpected output.

### Use Cases

* Displaying raw Mustache tags: When you want to show Mustache-style syntax without Vue interpreting it as data binding.

* Performance optimization: In situations where a template or part of a template doesn't contain any dynamic bindings, v-pre can prevent Vue from performing unnecessary compilation and diffing.

* Security: If you're rendering user-generated content that might contain Mustache-like syntax, v-pre can help prevent potential XSS vulnerabilities by treating the content as plain text.
Important points

* `v-pre´ affects the element it's applied to and all its descendants.
It's typically used in scenarios where you want to bypass Vue's template compilation for specific parts of your application.

## v-html & Escaping in general
by default vue only outputs escaped html, this directive allows to output unescaped html

### Important Considerations

* Security: Only use v-html with trusted content. Never use it to render user-generated content directly, as it can lead to Cross-Site Scripting (XSS) vulnerabilities.

* Alternatives: If possible, prefer using data binding with Mustache syntax ({{ }}) or computed properties to dynamically render content, as it's generally safer than v-html.

* Sanitization: If you must use v-html with potentially untrusted data, sanitize the input thoroughly using a trusted library like DOMPurify to remove any potentially harmful scripts or tags.

## Web workers
```html
<!DOCTYPE html>
<html>
<head>
  <title>WASM with Vue.js and Web Worker</title>
</head>
<body>
  <div id="app"></div>

  <script>
    // Create the Web Worker
    const worker = new Worker('my-worker.js');

    // Listen for messages from the worker
    worker.onmessage = (event) => {
      if (event.data.type === 'render') {
        document.getElementById('app').innerHTML = event.data.html;
      } else if (event.data.type === 'log') {
        console.log(event.data.message);
      }
    };

    // Send initial data to the worker
    worker.postMessage({ type: 'init' });
  </script>
</body>
</html>
```

**my-worker.js (Web Worker)**

```javascript
importScripts('vue.js', 'external-library.js'); // Import Vue and other libraries

onmessage = (event) => {
  if (event.data.type === 'init') {
    // Initialize Vue app
    const app = Vue.createApp({
      // ... your Vue app configuration
    });

    // Mount the app
    app.mount('#app');

    // Send initial rendered HTML to the main thread
    postMessage({ type: 'render', html: document.getElementById('app').innerHTML });
  } else if (event.data.type === 'update') {
    // Update Vue data or component state
    // ...

    // Re-render the relevant parts of the DOM
    // ...

    // Send updated HTML to the main thread
    postMessage({ type: 'render', html: /* updated HTML */ });
  }
};

// Function to call an external library
function callExternalLibrary(data) {
  // ... call the external library using its API
  externalLibrary.someFunction(data);
}
```

# Important Document for reactivity ideas

[signals RFC](https://github.com/tc39/proposal-signals)