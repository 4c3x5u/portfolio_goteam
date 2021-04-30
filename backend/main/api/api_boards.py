from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..models import Board, User
from ..validation.val_auth import \
    authenticate, authorize, not_authenticated_response, \
    not_authorized_response
from ..validation.val_team import validate_team_id
from ..validation.val_board import validate_board_id
from ..util import create_board


@api_view(['GET', 'POST', 'DELETE', 'PATCH'])
def boards(request):
    username = request.META.get('HTTP_AUTH_USER')
    token = request.META.get('HTTP_AUTH_TOKEN')

    user, authentication_response = authenticate(username, token)
    if authentication_response:
        return authentication_response

    # not in use – maintained for demonstration purposes
    if request.method == 'GET':
        if 'id' in request.query_params.keys():
            board_id = request.query_params.get('id')

            validation_response = validate_board_id(board_id)
            if validation_response:
                return validation_response

            try:
                board, = Board.objects.prefetch_related(
                    'user',
                    'column_set',
                    'column_set__task_set',
                    'column_set__task_set__subtask_set'
                ).get(id=board_id),
            except Board.DoesNotExist:
                return Response({
                    'board_id': ErrorDetail(string='Board not found.',
                                            code='not_found')
                }, 404)

            validation_response = validate_board_id(board_id)
            if validation_response:
                return validation_response

            if board.team.id != user.team.id:
                return not_authenticated_response

            try:
                board.user.get(username=user.username)
            except User.DoesNotExist:
                if not user.is_admin:
                    return not_authorized_response

            def column_mapper(column):
                return {
                    'id': column.id,
                    'order': column.order,
                    'tasks': column.task_set is not None and list(map(
                        lambda task: {
                            'id': task.id,
                            'title': task.title,
                            'description': task.description,
                            'order': task.order,
                            'user': task.user.username if task.user else '',
                            'subtasks': task.subtask_set is not None and list(map(
                                lambda subtask: {
                                    'id': subtask.id,
                                    'title': subtask.title,
                                    'order': subtask.order,
                                    'done': subtask.done
                                },
                                task.subtask_set.all()
                            ))
                        },
                        column.task_set.all()
                    ))
                }

            columns = list(map(column_mapper, board.column_set.all()))
            return Response({'id': board.id, 'columns': columns}, 200)

        if 'team_id' in request.query_params.keys():
            request_team_id = request.query_params.get('team_id')
            team, response = validate_team_id(request_team_id)
            if response:
                return response

            if team.id != user.team.id:
                return not_authenticated_response

            if user.is_admin:
                queryset = Board.objects.filter(team=team.id)
            else:
                queryset = Board.objects.filter(team=team.id, user=user)

            # create a board if none exists for the team and the user is admin
            if not queryset:
                if not authorize(username):
                    board, create_response = create_board(team.id, 'New Board')
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
        request_team_id = request.data.get('team_id')
        team, validation_response = validate_team_id(request_team_id)
        if validation_response:
            return validation_response

        if team.id != user.team.id:
            return not_authenticated_response

        board_name = request.data.get('name')

        board, create_response = create_board(team.id, board_name)
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

        validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response

        try:
            board = Board.objects.get(id=board_id)
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        if board.team.id != user.team.id:
            return not_authenticated_response

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
        validation_response = validate_board_id(board_id)
        if validation_response:
            return validation_response

        try:
            board = Board.objects.get(id=board_id)
        except Board.DoesNotExist:
            return Response({
                'board_id': ErrorDetail(string='Board not found.',
                                        code='not_found')
            }, 404)

        if board.team.id != user.team.id:
            return not_authenticated_response

        serializer = BoardSerializer(board, data=request.data, partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)

        serializer.save()
        return Response({
            'msg': 'Board updated successfuly.',
            'id': serializer.data['id'],
        }, 200)


