import React from 'react';
import PropTypes from "prop-types";
import {Input} from "reactstrap";
import makeGraphQLRequest from "../graphql";

class Pi extends Component {
    constructor(props) {
        super(props);

        this.changeDoor = this.changeDoor.bind(this);
    }


    changeDoor(e) {
        const query = `
        mutation ($id: Int!, $piId: Int!) {
            updateDoor(id: $id, piId: $piId) {
                id
            }
        }`;
        makeGraphQLRequest(query, {piId: this.props.id, id: e.target.value}, data => {
            if (data["data"] != null) {
                if (typeof this.props.onChange === "function") {
                    this.props.onChange();
                }
            }
        });
    }

    render() {
        let onlineText = null;
        if (this.props.online) {
            onlineText = <span className="text-success">Online</span>
        } else {
            onlineText = <span className="text-danger">Offline</span>
        }
        return (
            <tr>
                <th scope="row">{this.props.id}</th>
                <td>{this.props.mac}</td>
                <td>{onlineText}</td>
                <td>{this.props.lastSeen.toUTCString()}</td>
                <td>
                    <Input type="select" onChange={this.changeDoor} value={this.props.doorNum}>
                        {doors.map(door => (
                            <option key={door.id} value={door.id}>{door.number}</option>
                        ))}
                    </Input>
                </td>
            </tr>
        );
    }
}

Pi.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string.isRequired,
    online: PropTypes.bool.isRequired,
    doorNum: PropTypes.number,
    doors: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        number: PropTypes.number.isRequired
    })).isRequired,
    lastSeen: PropTypes.instanceOf(Date),
    onChange: PropTypes.func,
};

export default Pi;