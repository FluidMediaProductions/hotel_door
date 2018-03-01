import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Button, Modal, ModalHeader, ModalBody, ModalFooter, FormGroup, Form, Label, Input, Col} from 'reactstrap';
import makeGraphQLRequest from "../graphql";

class CreateDoor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: false
        };

        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
        this.create = this.create.bind(this);
    }

    show() {
        this.setState({
            modal: true
        });
    }

    hide() {
        this.setState({
            modal: false
        });
    }

    create() {
        const input = this.refs.input.refs.number;
        const query = `
        mutation ($number: Int!) {
            createDoor(number: $number) {
                id
            }
        }`;
        makeGraphQLRequest(query, {number: parseInt(input.value)}, data => {
            if (data["data"] != null) {
                if (typeof this.props.onChange === "function") {
                    this.props.onCreate();
                    this.hide();
                }
            }
        });
    }

    render() {
        return (
            <div>
                <Button color="success" onClick={this.show}>Add door</Button>
                <Modal isOpen={this.state.modal} toggle={this.hide}>
                    <ModalHeader toggle={this.hide}>Create door</ModalHeader>
                    <ModalBody>
                        <Form>
                            <FormGroup row>
                                <Label for="doorNumber" sm={4}>Number</Label>
                                <Col sm={8}>
                                    <Input type="number" id="doorNumber" innerRef="number" ref="input"/>
                                </Col>
                            </FormGroup>
                        </Form>
                    </ModalBody>
                    <ModalFooter>
                        <Button color="primary" onClick={this.create}>Do Something</Button>
                        <Button color="secondary" onClick={this.hide}>Cancel</Button>
                    </ModalFooter>
                </Modal>
            </div>
        )
    }
}

CreateDoor.propTypes = {
    onCreate: PropTypes.func
};

export default CreateDoor;