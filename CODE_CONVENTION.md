## Code convention

#### Read this before contribution or if you have a question how to write code in correct way

This project has pretty strict code of convention to prevent any misunderstandings, disputes in PRs and simplify code readability. All new proposals about
making changes to code convention should be discussed among contributors.

- **Styles (css)**
    - All styles (css) should be written as `css in js` by using `sx` prop (if it exists in component).
    - Styles shouldn't be located in `jsx`, the right way is to create constant above the component and to call it as `SX`, `style` or `class`


- **TypeScript**
    - All types should be created as `interface` in specific folder `app` (except props)
    - `Props` types should be created as `type` above the component
    - You should avoid creation of types inside the functions
    - All components should be written as function (you can create nested functions inside component if you want to separate logic and simplify code
      readability, but please try to avoid outer variables, especially if this variable can be prop)


- **GO**
    - in progress 
