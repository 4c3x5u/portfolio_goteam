import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import AppContext from '../../../AppContext';
import ColumnsAPI from '../../../api/ColumnsAPI';
import Column from './Column/Column';

import './board.sass';

const Board = ({ handleActivate }) => {
  const {
    user,
    activeBoard,
    loadBoard,
    setIsLoading,
    notify,
  } = useContext(AppContext);

  const handleOnDragEnd = async (result) => {
    try {
      setIsLoading(true);

      if (!user.isAdmin) {
        throw Error('Only the admin can do this.');
      }

      if (!result.destination) {
        throw Error('Destination not found.');
      }

      const source = activeBoard.columns.find((column) => (
        column.id.toString() === result.source.droppableId
      ));

      const [item] = source.tasks.splice(result.source.index, 1);

      const sourceTasks = source.tasks.map((task, index) => ({
        ...task,
        order: index,
      }));

      await ColumnsAPI.patch(source.id, sourceTasks);

      const destination = activeBoard.columns.find((column) => (
        column.id.toString() === result.destination.droppableId
      ));

      destination.tasks.splice(result.destination.index, 0, item);

      const destinationTasks = destination.tasks.map((task, index) => ({
        ...task,
        order: index,
      }));

      await ColumnsAPI.patch(destination.id, destinationTasks);

      await loadBoard(activeBoard.id);
    } catch (err) {
      notify(
        'Unable to update task.',
        `${err?.message || 'Server Error'}.`,
      );
      setIsLoading(false);
    }
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns && activeBoard.columns.map((column) => (
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
