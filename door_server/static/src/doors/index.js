import React, {Component} from 'react';
import {Col, Container, Row, Table} from 'reactstrap';
import makeGraphQLRequest from "../graphql";
import Door from "./Door";
import {paginationLength} from "../App";
import Pagination from "../Pagination";
import {getJWT} from "../auth";

class Doors extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
        this.state = {
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
                doorList(first: $first, offset: $offset) {
                    id,
                    pi {
                        mac
                    },
                    name
                }
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {token: getJWT(), first: paginationLength, offset: this.state.paginationOffset}, data => {
            if (data["data"]["auth"] != null) {
                let doors = [];
                for (const i in data["data"]["auth"]["doorList"]) {
                    const door = data["data"]["auth"]["doorList"][i];

                    doors.push({
                        id: door["id"],
                        name: door["name"],
                        mac: door["pi"]["mac"],
                    });
                }
                self.setState({
                    doors: doors
                })
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
        const nextDisabled = (this.state.doors.length <= paginationLength);
        return (
            <Container>
                <h1>Doors</h1>
                <Row>
                    <Col xs="12">
                        <Table hover>
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>Pi MAC</th>
                                <th>Door name</th>
                                <th>Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.doors.map(door => (
                                <Door key={door.id} id={door.id} piId={door.piId} mac={door.mac} name={door.name}/>
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

export default Doors;