from rest_framework.decorators import api_view
from rest_framework.response import Response
from uuid import UUID
from main.serializers import UserSerializer
from main.models import Team


@api_view(['POST'])
def register(request):
    password = request.data.get('password')
    password_confirmation = request.data.get('password_confirmation')
    if not password:
        return Response({'password': 'Password cannot be empty.'}, 400)
    if not password_confirmation:
        return Response({
            'password_confirmation': 'Password confirmation cannot be empty.'
        }, 400)
    if password != password_confirmation:
        return Response({
            'password_confirmation': 'Confirmation does not match the '
                                     'password.'
        }, 400)
    invite_code = request.data.get('invite_code')
    if invite_code:
        try:
            UUID(invite_code)
        except ValueError:
            return Response({'invite_code': 'Invalid invite code.'}, 400)
        try:
            team = Team.objects.get(invite_code=invite_code)
            is_admin = False
        except Team.DoesNotExist:
            return Response({'invite_code': "Team not found."}, 404)
    else:
        team = Team.objects.create()
        is_admin = True
    new_user = {'username': request.data.get('username'),
                'password': request.data.get('password'),
                'team': team.id,
                'is_admin': is_admin}
    serializer = UserSerializer(data=new_user)
    if serializer.is_valid():
        serializer.save()
        return Response(new_user, 201)
    is_admin and team.delete()
    return Response(serializer.errors, 400)

