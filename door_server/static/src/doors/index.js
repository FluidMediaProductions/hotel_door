import React, {Component} from 'react';
import {Col, Container, Row, Table} from 'reactstrap';
import makeGraphQLRequest from "../graphql";
import Door from "./Door";
import {paginationLength} from "../App";
import Pagination from "../Pagination";
import CreateDoor from "./CreateDoor";

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
        query ($first: Int!, $offset: Int!) {
            doorList(first: $first, offset: $offset) {
                id,
                pi {
                    id
                    mac
                },
                number
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {first: paginationLength, offset: this.state.paginationOffset}, data => {
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
                    <Col xs="12" className="text-right mb-3">
                        <CreateDoor onCreate={this.updateSate} />
                    </Col>
                </Row>
                <Row>
                    <Col xs="12">
                        <Table hover>
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>Pi ID</th>
                                <th>Pi MAC</th>
                                <th>Door number</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.doors.map(door => (
                                <Door key={door.id} id={door.id} piId={door.piId} mac={door.mac} number={door.number}/>
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