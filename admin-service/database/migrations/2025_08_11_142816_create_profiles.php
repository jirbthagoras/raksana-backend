<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('profiles', function (Blueprint $table) {
            $table->id();
            $table->foreignId("user_id")->references("id")->on("users");
            $table->bigInteger("current_exp")->default(0);
            $table->bigInteger("exp_needed")->default(100);
            $table->integer("level")->default(1);
            $table->bigInteger("points")->default(0);
            $table->string("profile_key")->default("profiles/Portrait_Placeholder.png");
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('profiles');
    }
};
