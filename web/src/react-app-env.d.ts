/// <reference types="react-scripts"/>
/*
These type declarations add support for importing resource files
such as bmp, gif, jpeg, jpg, png, webp, and svg. That means that
the following import will work as expected without errors:

`import logo from './logo.svg';`
*/

// Custom map types for EventSource.addEventListener()

interface EventSourceEventMap {
    "error": Event;
    "message": MessageEvent;
    "log": MessageEvent;
    "status": MessageEvent;
    "server": MessageEvent;
    "stream": MessageEvent;
    "open": Event;
}
