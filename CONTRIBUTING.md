# How to contribute

I'm really glad you're reading this, because we need volunteer developers to help this project come to fruition.

If you haven't already, you can write to us ([#ivory](https://postgresteam.slack.com/archives/C0304HP1ZSM) on postgres slack). We want you working on things you're excited about.

Here are some important resources:

  * [Enhancements](https://github.com/veegres/ivory/issues)
  * [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
  * [Supported versions](SECURITY.md)
  * [Setup Local Environment](docker/ivory-dev/README.md)
  * [Build Frontend](web/README.md)
  * [Build Backend](service/README.md)

## Testing

We don't have tests yet, so you need to test everything manually and good practice, before creating PR to make regression testing of all functionality. It is not required, but if you not sure or you change some critical part, it is better to spend some time doing this.

## Submitting changes

Please send a [GitHub Pull Request to Ivory](https://github.com/veegres/ivory/compare) with a clear list of what you've done (read more about [pull requests](http://help.github.com/pull-requests/)).

Always write a clear log message for your commits. One-line messages are fine for small changes, but bigger changes should look like this:

    $ git commit -m "A brief summary of the commit
    > 
    > A paragraph describing what changed and its impact."

## Coding conventions

Start reading our code and you'll get the hang of it. We optimize for readability:

- **Styles (css)**
    - All styles (css) should be written as `css in js` by using `sx` prop (if it exists in component).
    - Styles shouldn't be located in `jsx`, the right way is to create constant above the component and to call it as `SX`, `style` or `class`

- **TypeScript**
    - All types should be created as `interface` in specific folder `app` (except props)
    - `Props` types should be created as `type` above the component
    - You should avoid creation of types inside the functions
    - All components should be written as function (you can create nested functions inside component if you want to separate logic and simplify code
      readability, but please try to avoid outer variables, especially if this variable can be used as a prop)

- **GO**
    - in progress 


This project has pretty strict coding convention to prevent any misunderstandings, disputes in PRs and simplify code readability. All new proposals about
making changes to the coding convention should be discussed among contributors.


Thanks,
Andrei Sergeev
