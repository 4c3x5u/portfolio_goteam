from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..models import Board, Column, Task, Subtask
from ..util import (
    authenticate, authorize, validate_team_id, validate_board_id, create_board
)


@api_view(['GET', 'POST', 'DELETE', 'PATCH'])
def boards(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    if request.method == 'GET':
        if 'id' in request.query_params.keys():
            board_id = request.query_params.get('id')
            board, validation_response = validate_board_id(board_id)
            if validation_response:
                return validation_response

            columns = list(map(
                lambda column: {
                    'id': column.id,
                    'order': column.order,
                    'tasks': list(map(
                        lambda task: {
                            'id': task.id,
                            'title': task.title,
                            'description': task.description,
                            'order': task.order,
                            'subtasks': list(map(
                                lambda subtask: {
                                    'id': subtask.id,
                                    'title': subtask.title,
                                    'order': subtask.order,
                                    'done': subtask.done
                                },
                                Subtask.objects.filter(task_id=task.id)
                            ))
                        },
                        Task.objects.filter(column_id=column.id)
                    ))
                },
                Column.objects.filter(board_id=board.id)
            ))
            return Response({'id': board.id, 'columns': columns}, 200)

        if 'team_id' in request.query_params.keys():
            team_id = request.query_params.get('team_id')
            validation_response = validate_team_id(team_id)
            if validation_response:
                return validation_response

            # create a board if none exists for the team
            queryset = Board.objects.filter(team=team_id)
            if not queryset:
                if not authorize(username):
                    board, create_response = create_board(team_id, 'New Board')
                    if create_response:
                        return create_response

                    # return a list containing only the new board
                    return Response([{
                        'id': board.id, 'name': board.name
                    }], 201)

                return Response({
                    'team_id': ErrorDetail(string='Boards not found.',
                                           code='not_found')
                }, 404)

            return Response(list(map(
                lambda board_data: {
                    'id': board_data['id'],
                    'name': board_data['name']
                },
                BoardSerializer(queryset, many=True).data
            )), 200)

        return Response(list(map(
            lambda board_data: {
                'id': board_data['id'],
                'name': board_data['name']
            },
            BoardSerializer(Board.objects.all(), many=True).data
        )), 200)

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        # validate team_id
        team_id = request.data.get('team_id')
        validation_response = validate_team_id(team_id)
        if validation_response:
            return validation_response

        board_name = request.data.get('name')

        board, create_response = create_board(team_id, board_name)
        if create_response:
            return create_response

        # return success response
        return Response({
            'msg': 'Board creation successful.',
            'id': board.id,
        }, 201)

    if request.method == 'DELETE':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        board_id = request.query_params.get('id')

        board, validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response
        board.delete()

        return Response({
            'msg': 'Board deleted successfully.',
            'id': board_id,
        })

    if request.method == 'PATCH':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        board_id = request.query_params.get('id')
        board, validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response

        serializer = BoardSerializer(board, data=request.data, partial=True)
        if serializer.is_valid():
            serializer.save()
            return Response({
                'msg': 'Board updated successfuly.',
                'id': serializer.data['id'],
            }, 200)


