import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Button, Modal, ModalHeader, ModalBody, ModalFooter, FormGroup, Form, Label, Col} from 'reactstrap';
import makeGraphQLRequest from "../graphql";

class CreateDoor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: false
        };

        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
        this.save = this.save.bind(this);
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

    save() {
        const input = this.refs.input;
        const query = `
        mutation ($id: Int!, $number: Int!) {
            updateDoor(id: $id, number: $number) {
                number
            }
        }`;
        makeGraphQLRequest(query, {id: this.props.id, number: parseInt(input.value)}, data => {
            if (data["data"] != null) {
                if (typeof this.props.onSave === "function") {
                    this.props.onSave();
                }
                this.hide();
            }
        });
    }

    render() {
        return (
            <span>
                <Button color="primary" onClick={this.show}>
                    <i className="material-icons">edit</i>
                </Button>
                <Modal isOpen={this.state.modal} toggle={this.hide}>
                    <ModalHeader toggle={this.hide}>Edit door</ModalHeader>
                    <ModalBody>
                        <Form>
                            <FormGroup row>
                                <Label for="doorNumber" sm={4}>Number</Label>
                                <Col sm={8}>
                                    <input className="form-control" type="number" id="doorNumber"
                                           ref="input" defaultValue={this.props.number}/>
                                </Col>
                            </FormGroup>
                        </Form>
                    </ModalBody>
                    <ModalFooter>
                        <Button color="primary" onClick={this.save}>Save</Button>
                        <Button color="secondary" onClick={this.hide}>Cancel</Button>
                    </ModalFooter>
                </Modal>
            </span>
        )
    }
}

CreateDoor.propTypes = {
    id: PropTypes.number.isRequired,
    number: PropTypes.number.isRequired,
    onSave: PropTypes.func
};

export default CreateDoor;