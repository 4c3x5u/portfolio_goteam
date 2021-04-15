import React, { useContext } from 'react';
import PropTypes from 'prop-types';
// import axios from 'axios';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import Column from './Column/Column';
import AppContext from '../../../AppContext';

import './board.sass';

const Board = ({ handleActivate }) => {
  const { activeBoard } = useContext(AppContext);

  // TODO: API call to the database here inside a useEffect

  const handleOnDragEnd = (result) => {
    if (!result.destination) return;

    const source = activeBoard.columns.find(
      (column) => column.id.toString() === result.source.droppableId,
    );

    const [item] = source.tasks.splice(result.source.index, 1);

    const newSource = {
      ...source,
      tasks: source.tasks.map(
        (task, index) => ({ ...task, order: index }),
      ),
    };

    const destination = activeBoard.columns.find(
      (column) => column.id.toString() === result.destination.droppableId,
    );

    destination.tasks.splice(result.destination.index, 0, item);

    const newDestination = {
      ...destination,
      tasks: destination.tasks.map(
        (task, index) => ({ ...task, order: index }),
      ),
    };

    const newColumns = (
      activeBoard.columns.map((column) => {
        switch (column.id) {
          case destination.id: return newDestination;
          case source.id: return newSource;
          default: return column;
        }
      })
    );

    // TODO: handle
    console.log({ ...activeBoard, columns: newColumns });
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns.map((column) => (
          <Column
            id={column.id}
            name={column.order}
            tasks={column.tasks}
            handleActivate={handleActivate}
          />
        ))}
      </DragDropContext>
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
