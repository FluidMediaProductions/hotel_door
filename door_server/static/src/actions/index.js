import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Action from "./Action";

class Actions extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.state = {
            actions: []
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
            actionList {
                id
                pi {
                    mac
                    id
                }
                type
                success
                complete
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {}, data => {
            if (data["data"] != null) {
                let actions = [];
                for (const i in data["data"]["actionList"]) {
                    const action = data["data"]["actionList"][i];

                   actions.push({
                       id: action["id"],
                       type: action["type"],
                       mac: action["pi"]["mac"],
                       piId: action["pi"]["id"],
                       success: action["success"],
                       complete: action["complete"],
                    });
                }
                self.setState({
                    actions: actions
                });
            }
        });
    }

    render() {
        return (
            <div className="Doors container">
                <h1>Actions</h1>
                <div className="row">
                    <div className="col-12">
                        <table className="table table-hover">
                            <thead className="thead-light">
                            <tr>
                                <th scope="col">ID</th>
                                <th scope="col">Type</th>
                                <th scope="col">Pi ID</th>
                                <th scope="col">Pi MAC</th>
                                <th scope="col">Complete</th>
                                <th scope="col">Success</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.actions.map(action => (
                                <Action key={action.id} id={action.id} type={action.type} piId={action.piId} piMac={action.mac}
                                    success={action.success} complete={action.complete} />
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
}

export default Actions;