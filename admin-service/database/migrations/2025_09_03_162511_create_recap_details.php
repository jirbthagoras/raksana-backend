<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('recap_details', function (Blueprint $table) {
            $table->id();
            $table->foreignId("monthly_recap_id")->references("id")->on("recaps");$table->integer("challenges")->default(0);
            $table->integer("events")->default(0);
            $table->integer("quests")->default(0);
            $table->integer("treasures")->default(0);
            $table->integer("longest_streak")->default(0);
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('recap_details');
    }
};
