import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Button, Modal, ModalBody, ModalFooter, ModalHeader} from 'reactstrap';
import makeGraphQLRequest from "../graphql";
import {getJWT} from "../auth";

class DeletePi extends Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: false
        };

        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
        this.delete = this.delete.bind(this);
    }

    delete() {
        const query = `
        mutation ($token: String!, $id: Int!) {
            auth(token: $token) {
                deletePi(id: $id) {
                    deletedAt
                }
            }
        }`;
        makeGraphQLRequest(query, {id: this.props.id, token: getJWT()}, data => {
            if (data["data"]["auth"] != null) {
                if (typeof this.props.onDelete === "function") {
                    this.props.onDelete();
                }
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
                <Button color="danger" onClick={this.show} className="mr-2">
                    <i className="material-icons">delete</i>
                </Button>
                <Modal isOpen={this.state.modal} toggle={this.hide}>
                    <ModalHeader toggle={this.hide}>Delete pi</ModalHeader>
                    <ModalBody>
                        Are you sure you want to delete this pi?
                    </ModalBody>
                    <ModalFooter>
                        <Button color="danger" onClick={this.delete}>Confirm</Button>
                        <Button color="secondary" onClick={this.hide}>Cancel</Button>
                    </ModalFooter>
                </Modal>
            </span>
        )
    }
}

DeletePi.propTypes = {
    id: PropTypes.number.isRequired,
    onDelete: PropTypes.func
};

export default DeletePi;