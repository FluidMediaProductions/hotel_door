import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Pi from "./Pi";
import {paginationLength} from "../App";
import Pagination from '../Pagination';
import {Col, Container, Row, Table} from "reactstrap";
import {getJWT} from "../auth";

class Pis extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
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
        query ($token: String!, $first: Int!, $offset: Int!) {
            auth(token: $token) {
                piList(first: $first, offset: $offset) {
                    id,
                    mac,
                    online,
                    lastSeen
                    door {
                        id
                        name
                    }
                }
                doorList {
                    id
                    name
                }
            }
        }`;
        makeGraphQLRequest(query, {token: getJWT(), first: paginationLength, offset: this.state.paginationOffset}, data => {
            if (data["data"]["auth"] != null) {
                let pis = [];
                for (const i in data["data"]["auth"]["piList"]) {
                    const pi = data["data"]["auth"]["piList"][i];
                    let door = {
                        "name": null,
                        "id": null,
                    };
                    if (pi["door"] != null) {
                        door = pi["door"]
                    }

                    pis.push({
                        id: pi["id"],
                        mac: pi["mac"],
                        online: pi["online"],
                        lastSeen: new Date(pi["lastSeen"]),
                        doorName: door["name"],
                        doorId: door["id"]
                    });
                }
                let doors = [];
                for (const i in data["data"]["auth"]["doorList"]) {
                    const door = data["data"]["auth"]["doorList"][i];
                    doors.push({
                        id: door["id"],
                        name: door["name"]
                    })
                }
                this.setState({
                    pis: pis,
                    doors: doors
                });
            }
        });
    }

    nextPage(e) {
        e.preventDefault();
        this.setState((previousState) => ({
            paginationOffset: previousState.paginationOffset + paginationLength
        }), this.updateSate);
    }

    previousPage(e) {
        e.preventDefault();
        this.setState((previousState) => {
            let offset = previousState.paginationOffset - paginationLength;
            offset = (offset < 0) ? (0) : (offset);
            return {
                paginationOffset: offset
            }
        }, this.updateSate);
    }

    render() {
        const previousDisabled = (this.state.paginationOffset <= 0);
        const nextDisabled = (this.state.pis.length <= paginationLength);
        return (
            <Container>
                <h1>Pis</h1>
                <Row>
                    <Col xs="12">
                        <Table hover>
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>MAC</th>
                                <th>Online</th>
                                <th>Last Seen</th>
                                <th>Door Name</th>
                                <th>Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.pis.map(pi => (
                                <Pi key={pi.id} id={pi.id} mac={pi.mac} online={pi.online} lastSeen={pi.lastSeen}
                                    doorName={pi.doorName} doorId={pi.doorId} doors={this.state.doors} onChange={this.updateSate}/>
                            ))}
                            </tbody>
                        </Table>
                        <Pagination previousDisabled={previousDisabled} nextDisabled={nextDisabled}
                                    nextPage={this.nextPage} previousPage={this.previousPage}/>
                    </Col>
                </Row>
            </Container>
        );
    }
}

export default Pis;