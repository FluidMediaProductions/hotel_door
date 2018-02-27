import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Pi from "./Pi";

class Pis extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.state = {
            pis: []
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
            piList {
                id,
                mac,
                online,
                lastSeen
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {}, data => {
            if (data["data"] != null) {
                let pis = [];
                for (const i in data["data"]["piList"]) {
                    const pi = data["data"]["piList"][i];

                   pis.push({
                        id: pi["id"],
                        mac: pi["mac"],
                        online: pi["online"],
                        lastSeen: new Date(pi["lastSeen"])
                    });
                }
                self.setState({
                    pis: pis
                });
            }
        });
    }

    render() {
        return (
            <div className="Doors container">
                <h1>Pis</h1>
                <div className="row">
                    <div className="col-12">
                        <table className="table">
                            <thead>
                            <tr>
                                <th scope="col">ID</th>
                                <th scope="col">MAC</th>
                                <th scope="col">Online</th>
                                <th scope="col">Last Seen</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.pis.map(pi => (
                                <Pi key={pi.id} id={pi.id} mac={pi.mac} online={pi.online} lastSeen={pi.lastSeen} />
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}

export default Pis;