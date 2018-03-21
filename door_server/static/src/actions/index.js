import React, {Component} from 'react';
import makeGraphQLRequest from "../graphql";
import Action from "./Action";
import {paginationLength} from "../App";
import Pagination from "../Pagination";
import {Col, Container, Row, Table} from "reactstrap";
import {getJWT} from "../auth";

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
        query ($token: String!, $first: Int!, $offset: Int!) {
            auth(token: $token) {
                actionList(first: $first, offset: $offset) {
                    id
                    pi {
                        mac
                    }
                    type
                    success
                    complete
                }
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {first: paginationLength, offset: this.state.paginationOffset, token: getJWT()}, data => {
            if (data["data"]["auth"] != null) {
                let actions = [];
                for (const i in data["data"]["auth"]["actionList"]) {
                    const action = data["data"]["auth"]["actionList"][i];

                   actions.push({
                       id: action["id"],
                       type: action["type"],
                       mac: action["pi"]["mac"],
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
        const nextDisabled = (this.state.actions.length <= paginationLength);
        return (
            <Container>
                <h1>Actions</h1>
                <Row>
                    <Col xs="12">
                        <Table hover>
                            <thead>
                            <tr>
                                <th>ID</th>
                                <th>Type</th>
                                <th>Pi MAC</th>
                                <th>Complete</th>
                                <th>Success</th>
                            </tr>
                            </thead>
                            <tbody>
                            {this.state.actions.map(action => (
                                <Action key={action.id} id={action.id} type={action.type} piMac={action.mac}
                                    success={action.success} complete={action.complete} />
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

export default Actions;