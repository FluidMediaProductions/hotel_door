import React, {Component} from 'react';
import PropTypes from "prop-types";
import makeGraphQLRequest from "../graphql";

class Pi extends Component {
    constructor(props) {
        super(props);

        this.updateSate = this.updateSate.bind(this);
        this.changeDoor = this.changeDoor.bind(this);
        this.state = {
            doorId: null,
            doors: [],
        }
    }

    componentDidMount() {
        this.updateSate();
        this.timer = setInterval(this.updateSate, 1000);
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    updateSate() {
        const query = `
        query ($id: Int!) {
            pi(id: $id) {
                door {
                    id
                }
            }
            doorList {
                id
                number
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {id: this.props.id}, data => {
            if (data["data"] != null) {
                const pi = data["data"]["pi"];

                let doors = [];
                for (const i in data["data"]["doorList"]) {
                    const door = data["data"]["doorList"][i];
                    doors.push({
                        id: door["id"],
                        number: door["number"]
                    })
                }

                self.setState({
                    doorId: pi["door"]["id"],
                    doors: doors
                })
            }
        });
    }

    changeDoor(e) {
        const query = `
        mutation ($id: Int!, $piId: Int!) {
            updateDoor(id: $id, piId: $piId) {
                id
            }
        }`;
        const self = this;
        makeGraphQLRequest(query, {piId: this.props.id, id: e.target.value}, data => {
            if (data["data"] != null) {
                const door = data["data"]["updateDoor"];

                self.setState({
                    doorId: door["id"]
                })
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
                    <select className="custom-select" onChange={this.changeDoor} value={this.state.doorId} ref="doorSelect">
                        {this.state.doors.map(door => (
                           <option key={door.id} value={door.id}>{door.number}</option>
                        ))}
                    </select>
                </td>
            </tr>
        );
    }
}

Pi.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string.isRequired,
    online: PropTypes.bool.isRequired,
    lastSeen: PropTypes.instanceOf(Date),
};

export default Pi;