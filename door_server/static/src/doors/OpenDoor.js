import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Button, Modal, ModalBody, ModalFooter, ModalHeader} from 'reactstrap';
import makeGraphQLRequest from "../graphql";
import {getJWT} from "../auth";

class OpenDoor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: false
        };

        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
        this.open = this.open.bind(this);
    }

    open() {
        const query = `
        mutation ($token: String!, $id: Int!) {
            auth(token: $token) {
                openDoor(id: $id)
            }
        }`;
        makeGraphQLRequest(query, {id: this.props.id, token: getJWT()}, data => {
            if (data["data"]["auth"] != null) {
                if (typeof this.props.onOpen === "function") {
                    this.props.onOpen();
                }
                this.hide();
            }
        });
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

    render() {
        return (
            <span>
                <Button color="success" onClick={this.show} className="mr-2">
                    <i className="material-icons">lock_open</i>
                </Button>
                <Modal isOpen={this.state.modal} toggle={this.hide}>
                    <ModalHeader toggle={this.hide}>Open door</ModalHeader>
                    <ModalBody>
                        Are you sure you want to open the door?
                    </ModalBody>
                    <ModalFooter>
                        <Button color="primary" onClick={this.open}>Open</Button>
                        <Button color="secondary" onClick={this.hide}>Cancel</Button>
                    </ModalFooter>
                </Modal>
            </span>
        )
    }
}

OpenDoor.propTypes = {
    id: PropTypes.number.isRequired,
    onOpen: PropTypes.func
};

export default OpenDoor;