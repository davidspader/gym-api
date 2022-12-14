<?php

namespace Tests\Feature\API\user;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Testing\Fluent\AssertableJson;
use Tests\TestCase;

class LoginControllerTest extends TestCase
{
    use RefreshDatabase;

    public function test_login_endpoint()
    {
        $user = User::factory(1)->createOne();

        $response = $this->postJson('/api/auth/login', [
            'email' => $user->email,
            'password' => 'password',
        ]);

        $response->assertStatus(200);

        $response->assertJson(function (AssertableJson $json) use($user){

            $json->hasAll(['data.user', 'data.token']);

            $json->whereAll([
                'data.user.name' => $user['name'],
                'data.user.email' => $user['email']
            ])->etc();
        });
    }

    public function test_fail_login_endpoint()
    {
        $user = User::factory(1)->createOne();

        $response = $this->postJson('/api/auth/login', [$user]);

        $response->assertStatus(401);

        $response->assertJson(function (AssertableJson $json) use($user){

            $json->hasAll(['error', 'error.message']);

            $json->whereAll([
                'error.message' => 'Invalid credentials',
            ]);
        });
    }

    public function test_logout_without_token()
    {
        $response = $this->postJson('/api/auth/logout');

        $response->assertStatus(401);

        $response->assertJson(function (AssertableJson $json) {

            $json->hasAll('message');

            $json->whereAll([
                'message' => 'Unauthenticated.',
            ]);
        });
    }
}
