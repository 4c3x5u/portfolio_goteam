from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..serializers.ser_board import BoardSerializer
from ..models import Board, User, Team
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

            if board.team_id != user.team_id:
                return not_authenticated_response

            try:
                board.user.get(username=user.username)
            except User.DoesNotExist:
                if not user.is_admin:
                    return not_authorized_response

            return Response({
                'id': board.id,
                'columns': [
                    {
                        'id': column.id,
                        'order': column.order,
                        'tasks': column.task_set is not None and [
                            {
                                'id': task.id,
                                'title': task.title,
                                'description': task.description,
                                'order': task.order,
                                'user': task.user.username or '',
                                'subtasks': task.subtask_set is not None and [
                                    {
                                        'id': subtask.id,
                                        'title': subtask.title,
                                        'order': subtask.order,
                                        'done': subtask.done
                                    } for subtask in task.subtask_set.all()
                                ]
                            } for task in column.task_set.all()
                        ]
                    } for column in board.column_set.all()
                ]
            }, 200)

        if 'team_id' in request.query_params.keys():
            team_id = request.query_params.get('team_id')
            response = validate_team_id(team_id)
            if response:
                return response

            try:
                team = Team.objects.prefetch_related(
                    'board_set',
                    'user_set',
                ).get(id=team_id)
            except Team.DoesNotExist:
                return Response({
                    'team_id': ErrorDetail(string='Team not found.',
                                           code='not_found')
                }, 404)

            if team.id != user.team_id:
                return not_authenticated_response

            if user.is_admin:
                board_list = team.board_set.all()
            else:
                board_list = team.board_set.filter(user=user)

            # create a board if none exists for the team and the user is admin
            if not board_list:
                if not authorize(username):
                    board, error_response = create_board(
                        name='New Board',
                        team_id=team.id,
                        team_admin=team.user_set.get(is_admin=True)
                    )
                    if error_response:
                        return error_response

                    # return a list containing only the new board
                    return Response([
                        {'id': board.id, 'name': board.name}
                    ], 201)

                return Response({
                    'team_id': ErrorDetail(string='Boards not found.',
                                           code='not_found')
                }, 404)

            return Response([
                {'id': board.id, 'name': board.name} for board in board_list
            ], 200)

        board_list = Board.objects.all()
        return Response([
            {'id': board.id, 'name': board.name} for board in board_list
        ], 200)

    if request.method == 'POST':
        authorization_response = authorize(username)
        if authorization_response:
            return authorization_response

        # validate team_id
        team_id = request.data.get('team_id')
        validation_response = validate_team_id(team_id)
        if validation_response:
            return validation_response

        try:
            team = Team.objects.prefetch_related('user_set').get(id=team_id)
        except Team.DoesNotExist:
            return Response({
                'team_id': ErrorDetail(string='Team not found.',
                                       code='not_found')
            }, 404)

        if team.id != user.team_id:
            return not_authenticated_response

        board_name = request.data.get('name')

        board, error_response = create_board(
            name=board_name,
            team_id=team.id,
            team_admin=team.user_set.get(username=username)
        )
        if error_response:
            return error_response

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

        if board.team_id != user.team_id:
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

        if board.team_id != user.team_id:
            return not_authenticated_response

        serializer = BoardSerializer(board, data=request.data, partial=True)
        if not serializer.is_valid():
            return Response(serializer.errors, 400)
        serializer.save()

        return Response({
            'msg': 'Board updated successfuly.',
            'id': serializer.data['id'],
        }, 200)


