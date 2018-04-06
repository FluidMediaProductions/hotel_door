import React, {Component} from 'react';
import PropTypes from "prop-types";
import {Input} from "reactstrap";
import makeGraphQLRequest from "../graphql";
import DeletePi from "./DeletePi";
import {getJWT} from "../auth";

class Pi extends Component {
    constructor(props) {
        super(props);

        this.state = {
            modal: false
        };

        this.changeDoor = this.changeDoor.bind(this);
        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
    }


    changeDoor(e) {
        const query = `
        mutation ($token: String!, $id: Int!, $piId: Int!) {
            auth(token: $token) {
                updateDoor(id: $id, piId: $piId) {
                    id
                }
            }
        }`;
        makeGraphQLRequest(query, {piId: this.props.id, id: e.target.value, token: getJWT()}, data => {
            if (data["data"]["auth"] != null) {
                if (typeof this.props.onChange === "function") {
                    this.props.onChange();
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
                    <Input type="select" onChange={this.changeDoor} value={this.props.doorId}>
                        <option value="">-</option>
                        {this.props.doors.map(door => (
                            <option key={door.id} value={door.id}>{door.name}</option>
                        ))}
                    </Input>
                </td>
                <td>
                    <DeletePi id={this.props.id} onDelete={this.props.onChange}/>
                </td>
            </tr>
        );
    }
}

Pi.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string.isRequired,
    online: PropTypes.bool.isRequired,
    doorName: PropTypes.string,
    doorId: PropTypes.number,
    doors: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired
    })).isRequired,
    lastSeen: PropTypes.instanceOf(Date),
    onChange: PropTypes.func,
};

export default Pi;