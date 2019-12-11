# Mediator Examples

This is the repository containing examples from [the mediator game engine](https://github.com/celskeggs/mediator/).
It is a separate repository, because it involves examples ported from other peoples' code.

## License

All of the code I've written here is licensed under the MIT license.

Any code written by others is licensed under the appropriate license; please refer to the following links for more information:

 * yourfirstworld: http://www.byond.com/developer/Dantom/YourFirstWorld

# Trying it out

The easiest way to try this is to 'go get' the package.

    $ go get github.com/celskeggs/mediator-examples/yourfirstworld/
    $ yourfirstworld -core "$GOPATH/src/github.com/celskeggs/mediator/resources/" \
                     -resources "$GOPATH/src/github.com/celskeggs/mediator-examples/yourfirstworld/resources/" \
                     -map "$GOPATH/src/github.com/celskeggs/mediator-examples/yourfirstworld/map.dmm"

To run the example, you need to point it at the three resource locations shown above, which represent:

 * "core": the webclient resources, such as the HTML, JS, and CSS code for the webclient
 * "resources": the game-specified resources, such as icons
 * "map": the location of the BYOND DMM map

## Disclaimer

This is incredibly untested, and the examples don't fully work. Don't expect much from this.
