import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Action from "./Action";
import {paginationLength} from "../App";

class Actions extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
        this.state = {
            actions: [],
            paginationOffset: 0
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
        query ($first: Int!, $offset: Int!) {
            actionList(first: $first, offset: $offset) {
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
        makeGraphQLRequest(query, {first: paginationLength, offset: this.state.paginationOffset}, data => {
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

    nextPage(e) {
        e.preventDefault();
        this.setState((previousState) => ({
            paginationOffset: previousState.paginationOffset+paginationLength
        }), this.updateSate);
    }

    previousPage(e) {
        e.preventDefault();
        this.setState((previousState) => {
            let offset = previousState.paginationOffset-paginationLength;
            offset = (offset < 0)?(0):(offset);
            return {
                paginationOffset: offset
            }
        }, this.updateSate);
    }

    render() {
        const previousDisabled = (this.state.paginationOffset <= 0);
        const nextDisabled = (this.state.actions.length === 0);
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
                        <nav>
                            <ul className="pagination justify-content-center">
                                <li className={"page-item"+(previousDisabled?(" disabled"):(""))}>
                                    <a className="page-link" href="" onClick={this.previousPage}>Previous</a>
                                </li>
                                <li className={"page-item"+(nextDisabled?(" disabled"):(""))}>
                                    <a className="page-link" href="" onClick={this.nextPage}>Next</a>
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>
            </div>
        );
    }
}

export default Actions;