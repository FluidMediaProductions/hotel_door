import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Door from "./Door";

class Doors extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.state = {
            doors: []
        }
    }

    componentDidMount() {
        this.updateSate();
        this.timer = setInterval(this.updateSate, 1000);
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    updateSate() {
        const query = `
        query {
            doorList {
                id,
                pi {
                    id
                    mac
                },
                number
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {}, data => {
            if (data["data"] != null) {
                let doors = [];
                for (const i in data["data"]["doorList"]) {
                    const door = data["data"]["doorList"][i];

                    doors.push({
                        id: door["id"],
                        number: door["number"],
                        mac: door["pi"]["mac"],
                        piId: door["pi"]["id"],
                    });
                }
                self.setState({
                    doors: doors
                })
            }
        });
    }

    render() {
        return (
            <div className="Doors container">
                <h1>Doors</h1>
                <div className="row">
                    <div className="col-12">
                        <table className="table">
                            <thead>
                            <tr>
                                <th scope="col">ID</th>
                                <th scope="col">Pi ID</th>
                                <th scope="col">Pi MAC</th>
                                <th scope="col">Door number</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.doors.map(door => (
                                <Door key={door.id} id={door.id} piId={door.piId} mac={door.mac} number={door.number}/>
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}

export default Doors;