import React, {Component} from 'react';
import PropTypes from "prop-types";
import {Button} from "reactstrap";
import makeGraphQLRequest from "../graphql";
import EditDoor from "./EditDoor";

class Door extends Component {
    constructor(props) {
        super(props);

        this.delete = this.delete.bind(this);
    }

    delete() {
        const query = `
        mutation ($id: Int!) {
            deleteDoor(id: $id) {
                deletedAt
            }
        }`;
        makeGraphQLRequest(query, {id: this.props.id}, data => {
            if (data["data"] != null) {
                if (typeof this.props.onChange === "function") {
                    this.props.onUpdate();
                }
            }
        });
    }

    render() {
        return (
            <tr>
                <th scope="row">{this.props.id}</th>
                <td>{this.props.mac}</td>
                <td>{this.props.number}</td>
                <td>
                    <Button color="danger" onClick={this.delete} className="mr-2">
                        <i className="material-icons">delete</i>
                    </Button>
                    <EditDoor id={this.props.id} number={this.props.number} onSave={this.props.onUpdate}/>
                </td>
            </tr>
        )
    }
}

Door.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string,
    number: PropTypes.number.isRequired,
    onUpdate: PropTypes.func
};

export default Door;