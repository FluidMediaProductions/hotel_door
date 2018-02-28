import React from 'react';
import PropTypes from "prop-types";

const Pi = ({id, mac, online, doorNum, lastSeen, doors, onChange}) => {
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
            <td>
                <select className="custom-select" data-id={id} onChange={onChange}
                        value={doorNum}>
                    {doors.map(door => (
                       <option key={door.id} value={door.id}>{door.number}</option>
                    ))}
                </select>
            </td>
        </tr>
    );
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