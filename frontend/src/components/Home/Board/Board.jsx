import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';

import Column from './Column/Column';
import { columnOrder } from './Column/columnOrder';

import './board.sass';

const Board = ({ handleActivate }) => {
  const [columns] = useState([
    { id: 0, order: columnOrder.INBOX },
    { id: 1, order: columnOrder.READY },
    { id: 2, order: columnOrder.GO },
    { id: 3, order: columnOrder.DONE },
  ]);

  // TODO: API call to the database here inside a useEffect
  // useEffect(() => (
  //   setColumns([
  //   ])
  // ), []);

  const getColumnId = (order) => (
    columns.find((column) => column.order === order).id
  );

  return (
    <Row id="Board">
      <Column
        id={getColumnId(columnOrder.INBOX)}
        name={columnOrder.INBOX}
        handleActivate={handleActivate}
      />

      <Column id={getColumnId(columnOrder.READY)} name={columnOrder.READY} />

      <Column id={getColumnId(columnOrder.GO)} name={columnOrder.GO} />

      <Column id={getColumnId(columnOrder.DONE)} name={columnOrder.DONE} />
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
