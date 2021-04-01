from rest_framework import serializers
from ..models import User, Team, Board


class BoardSerializer(serializers.ModelSerializer):
    username = serializers.CharField(min_length=5,
                                     max_length=35,
                                     required=False)
    team_id = serializers.IntegerField(
        error_messages={'null': 'Team ID cannot be empty.'}
    )

    class Meta:
        model = Board
        fields = ('id', 'team_id', 'username')

    def update(self, instance, validated_data):
        raise Exception('Update not allowed on this endpoint.')

    def validate(self, data):
        try:
            Team.objects.get(id=data.get('team_id'))
        except Team.DoesNotExist:
            raise serializers.ValidationError({
                'team_id': 'Invalid team ID.',
            })

        return super().validate(data)

    def create(self, validated_data):
        username = validated_data.get('username')
        if not username:
            raise serializers.ValidationError({
                'username': "Usernme cannot be empty."
            }, 'null')
        try:
            user = User.objects.get(username=validated_data.get('username'))
        except User.DoesNotExist:
            raise serializers.ValidationError({
                'username': "Invalid username."
            }, 'invalid')
        if not user.is_admin:
            raise serializers.ValidationError({
                'is_admin': 'Only the team admin can create a board.',
            }, 'not_authorized')
        return Board.objects.create(team=user.team)
