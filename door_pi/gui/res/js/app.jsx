/** @jsx preact.h */

class App extends preact.Component {
    render() {
        let display = null;
        if (!app.data.opening) {
            display = (
                <div>
                    <h2>Room number</h2>
                    <h1 className="display-1">{app.data.name}</h1>
                </div>
            );
        } else {
            display = (
                <div>
                    <h1 className="display-2">Opening door...</h1>
                </div>
            )
        }
        return (
            <div class="container h-100">
                <div class="row align-items-center h-100">
                    <div className="col-12 text-center">
                        {display}
                    </div>
                </div>
            </div>
        );
    }
}

const render = () =>
    preact.render(<App />, document.getElementById('app'), document.getElementById('app').lastElementChild);

app.render = render;

render();