import React, {Component} from 'react';
import PropTypes from "prop-types";
import EditDoor from "./EditDoor";
import DeleteDoor from "./DeleteDoor";
import OpenDoor from "./OpenDoor";

class Door extends Component {
    render() {
        return (
            <tr>
                <th scope="row">{this.props.id}</th>
                <td>{this.props.mac}</td>
                <td>{this.props.number}</td>
                <td>
                    <DeleteDoor id={this.props.id} onDelete={this.props.onUpdate} />
                    <EditDoor id={this.props.id} number={this.props.number} onSave={this.props.onUpdate} />
                    <OpenDoor id={this.props.id} />
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