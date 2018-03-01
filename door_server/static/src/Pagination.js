import React from 'react';
import PropTypes from 'prop-types';

const Pagination = ({previousDisabled, nextDisabled, previousPage, nextPage}) => (
    <nav>
        <ul className="pagination justify-content-center">
            <li className={"page-item"+(previousDisabled?(" disabled"):(""))}>
                <a className="page-link" href="" onClick={previousPage}>Previous</a>
            </li>
            <li className={"page-item"+(nextDisabled?(" disabled"):(""))}>
                <a className="page-link" href="" onClick={nextPage}>Next</a>
            </li>
        </ul>
    </nav>
);

Pagination.propTypes = {
    previousDisabled: PropTypes.bool.isRequired,
    nextDisabled: PropTypes.bool.isRequired,
    previousPage: PropTypes.func.isRequired,
    nextPage: PropTypes.func.isRequired,
};

export default Pagination;