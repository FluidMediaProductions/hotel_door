import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Pi from "./Pi";
import {paginationLength} from "../App";

class Pis extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
        this.changeDoor = this.changeDoor.bind(this);
        this.state = {
            pis: [],
            doors: [],
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
            piList(first: $first, offset: $offset) {
                id,
                mac,
                online,
                lastSeen
                door {
                    number
                }
            }
            doorList {
                id
                number
            }
        }`;
        makeGraphQLRequest(query, {first: paginationLength, offset: this.state.paginationOffset}, data => {
            if (data["data"] != null) {
                let pis = [];
                for (const i in data["data"]["piList"]) {
                    const pi = data["data"]["piList"][i];

                   pis.push({
                        id: pi["id"],
                        mac: pi["mac"],
                        online: pi["online"],
                        lastSeen: new Date(pi["lastSeen"]),
                        doorNum: pi["door"]["number"]
                    });
                }
                let doors = [];
                for (const i in data["data"]["doorList"]) {
                    const door = data["data"]["doorList"][i];
                    doors.push({
                        id: door["id"],
                        number: door["number"]
                    })
                }
                this.setState({
                    pis: pis,
                    doors: doors
                });
            }
        });
    }

    changeDoor(e) {
        const query = `
        mutation ($id: Int!, $piId: Int!) {
            updateDoor(id: $id, piId: $piId) {
                id
            }
        }`;
        makeGraphQLRequest(query, {piId: e.target.dataset.id, id: e.target.value}, data => {
            if (data["data"] != null) {
                this.updateSate();
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
        const nextDisabled = (this.state.pis.length <= paginationLength);
        return (
            <div className="Doors container">
                <h1>Pis</h1>
                <div className="row">
                    <div className="col-12">
                        <table className="table table-hover">
                            <thead className="thead-light">
                            <tr>
                                <th scope="col">ID</th>
                                <th scope="col">MAC</th>
                                <th scope="col">Online</th>
                                <th scope="col">Last Seen</th>
                                <th scope="col">Door Number</th>
                                <th scope="col">Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.pis.map(pi => (
                                <Pi key={pi.id} id={pi.id} mac={pi.mac} online={pi.online} lastSeen={pi.lastSeen}
                                    doorNum={pi.doorNum} doors={this.state.doors} onChange={this.changeDoor} />
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

export default Pis;