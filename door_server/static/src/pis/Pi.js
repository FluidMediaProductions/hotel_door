import React, {Component} from 'react';
import PropTypes from "prop-types";
import {Input} from "reactstrap";
import makeGraphQLRequest from "../graphql";
import DeletePi from "./DeletePi";

class Pi extends Component {
    constructor(props) {
        super(props);

        this.state = {
            modal: false
        };

        this.changeDoor = this.changeDoor.bind(this);
        this.delete = this.delete.bind(this);
        this.show = this.show.bind(this);
        this.hide = this.hide.bind(this);
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

    delete() {
        const query = `
        mutation ($id: Int!) {
            deletePi(id: $id) {
                deletedAt
            }
        }`;
        makeGraphQLRequest(query, {id: this.props.id}, data => {
            if (data["data"] != null) {
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
                    <Input type="select" onChange={this.changeDoor} value={this.props.doorNum}>
                        <option value="">-</option>
                        {this.props.doors.map(door => (
                            <option key={door.id} value={door.id}>{door.number}</option>
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
    doorNum: PropTypes.number,
    doors: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        number: PropTypes.number.isRequired
    })).isRequired,
    lastSeen: PropTypes.instanceOf(Date),
    onChange: PropTypes.func,
};

export default Pi;