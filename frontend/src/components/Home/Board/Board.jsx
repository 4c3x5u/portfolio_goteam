import React, { useContext } from 'react';
import PropTypes from 'prop-types';
// import axios from 'axios';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import Column from './Column/Column';
import AppContext from '../../../AppContext';
import { patchColumn } from '../../../misc/api';

import './board.sass';

const Board = ({ handleActivate }) => {
  const { activeBoard, loadBoard, setIsLoading } = useContext(AppContext);

  // TODO: API call to the database here inside a useEffect

  const handleOnDragEnd = async (result) => {
    setIsLoading(true);

    if (!result.destination) return;

    const source = activeBoard.columns.find(
      (column) => column.id.toString() === result.source.droppableId,
    );

    const [item] = source.tasks.splice(result.source.index, 1);

    const sourceTasks = source.tasks.map(
      (task, index) => ({ ...task, order: index }),
    );

    try {
      await patchColumn(source.id, sourceTasks);
    } catch (err) {
      console.error(err);
    }

    const destination = activeBoard.columns.find(
      (column) => column.id.toString() === result.destination.droppableId,
    );

    destination.tasks.splice(result.destination.index, 0, item);

    const destinationTasks = destination.tasks.map(
      (task, index) => ({ ...task, order: index }),
    );

    try {
      await patchColumn(destination.id, destinationTasks);
    } catch (err) {
      console.error(err);
    }

    await loadBoard(activeBoard.id);
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns.map((column) => (
          <Column
            id={column.id}
            order={column.order}
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
