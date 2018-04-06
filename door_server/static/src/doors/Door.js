import React, {Component} from 'react';
import PropTypes from "prop-types";
import OpenDoor from "./OpenDoor";

class Door extends Component {
    render() {
        return (
            <tr>
                <th scope="row">{this.props.id}</th>
                <td>{this.props.mac}</td>
                <td>{this.props.name}</td>
                <td>
                    <OpenDoor id={this.props.id} />
                </td>
            </tr>
        )
    }
}

Door.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string,
    name: PropTypes.string.isRequired
};

export default Door;