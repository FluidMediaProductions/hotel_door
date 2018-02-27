import React from 'react';
import PropTypes from "prop-types";

const Pi = ({id, mac, online, lastSeen}) => {
    let onlineText = null;
    if (online) {
        onlineText = <span className="text-success">Online</span>
    } else {
        onlineText = <span className="text-danger">Offline</span>
    }
    return (
        <tr>
            <th scope="row">{id}</th>
            <td>{mac}</td>
            <td>{onlineText}</td>
            <td>{lastSeen.toUTCString()}</td>
        </tr>
    );
}

Pi.propTypes = {
    id: PropTypes.number.isRequired,
    mac: PropTypes.string.isRequired,
    online: PropTypes.bool.isRequired,
    lastSeen: PropTypes.instanceOf(Date),
};

export default Pi;